import React, { useEffect, useRef, useState, useMemo, useCallback } from 'react';
import * as d3 from 'd3';
import { TransformWrapper, TransformComponent } from 'react-zoom-pan-pinch';
import cloud from 'd3-cloud';
import { debounce } from 'lodash-es';
import { WordCloudItem, WordCloudConfig } from '@/lib/types/wordcloud';

interface WordCloudCanvasProps {
  data: WordCloudItem[];
  config?: Partial<WordCloudConfig>;
  onWordClick?: (word: WordCloudItem) => void;
}

const defaultConfig: WordCloudConfig = {
  minFontSize: 12,
  maxFontSize: 64,
  colors: {
    background: '#ffffff',
    wordColors: [
      '#2563eb', '#3b82f6', '#60a5fa', '#93c5fd', '#bfdbfe',
      '#1d4ed8', '#2563eb', '#3b82f6', '#60a5fa', '#93c5fd'
    ],
    highlightColor: '#1e40af'
  },
  layout: {
    padding: 3,
    spiral: 'archimedean',
    rotation: {
      angles: [0, 0, 0, 90], // 0度を多めに
      random: true
    }
  }
};

// 単語の重要度を計算するためのヘルパー関数
const calculateWordImportance = (word: string, count: number): number => {
  // 単語の長さによる重み付け（短すぎる単語は重要度を下げる）
  const lengthWeight = word.length < 2 ? 0.5 : word.length < 3 ? 0.8 : 1;
  
  // 出現回数による重み付け
  const countWeight = Math.log10(count + 1);
  
  return lengthWeight * countWeight;
};

// 不要な単語をフィルタリングするための関数
const shouldIncludeWord = (word: string): boolean => {
  // ボットメッセージや定型文を除外
  const botPatterns = [
    'リマインダー',
    'さんがチャンネルに参加しました',
    'http',
    'https',
    '議事録はこちら',
    'google',
    'doc',
    'Figma',
    'file'
  ];
  
  if (botPatterns.some(pattern => word.includes(pattern))) return false;

  // 特定の文字や記号のみの単語を除外
  const invalidPattern = /^[!-/:-@[-`{-~！-／：-＠［-｀｛-～、-〜"'・]+$/;
  if (invalidPattern.test(word)) return false;

  // 1文字の記号を除外
  if (word.length === 1 && /[!-/:-@[-`{-~！-／：-＠［-｀｛-～、-〜"'・]/.test(word)) return false;

  return true;
};

const processWordCloudData = (data: WordCloudItem[]): WordCloudItem[] => {
  // より積極的なフィルタリングと重み付け
  const emphasisWords = {
    '始めます！': 2.0,    // より強調
    'アジェンダ': 1.8,    // より強調
    '定例': 1.6,         // 重要な会議関連ワード
    '議題': 1.6,
    'スキップ': 1.4,
    '確認': 1.3,
    'よろしく': 1.3,
    'お願い': 1.3
  };

  // 除外するワードを追加
  const excludeWords = new Set([
    'cc', 'CC', '様', 'さん', 'こと', 'ため', 'http', 'https',
    'の', 'に', 'は', 'を', 'が', 'と', 'です', 'ます', 'した',
    'リマインダー', '議事録', 'Figma', 'google', 'doc'
  ]);

  return data
    .filter(item => {
      const text = item.text.trim();
      return text.length >= 2 && // 2文字以上
             !excludeWords.has(text) &&
             shouldIncludeWord(text);
    })
    .map(item => {
      const text = item.text.trim();
      const lengthMultiplier = text.length >= 3 && text.length <= 4 ? 1.5 : 1.0;
      const emphasisMultiplier = emphasisWords[text] || 1.0;
      
      // より鮮やかな配色
      const colors = [
        '#FF4B00', // ビビッドな赤
        '#005AFF', // 鮮やかな青
        '#00B06B', // 明るい緑
        '#FFB400', // オレンジ
        '#9900FF'  // 紫
      ];

      return {
        ...item,
        text,
        count: Math.floor(item.count * lengthMultiplier * emphasisMultiplier),
        color: colors[Math.floor(Math.random() * colors.length)],
        rotate: Math.random() > 0.5 ? 0 : 90, // ランダムな回転
        fontSize: calculateFontSize(item.count * lengthMultiplier * emphasisMultiplier)
      };
    })
    .sort((a, b) => b.count - a.count)
    .slice(0, 50); // より少ない単語数で見やすく
};

const calculateFontSize = (count: number): number => {
  // より大きなフォントサイズの差
  const minSize = 24;  // 最小フォントサイズを大きく
  const maxSize = 72;  // 最大フォントサイズをより大きく
  return Math.max(minSize, Math.min(maxSize, Math.floor(Math.log(count) * 12)));
};

export const WordCloudCanvas: React.FC<WordCloudCanvasProps> = ({
  data,
  config: userConfig,
  onWordClick
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const svgRef = useRef<SVGSVGElement>(null);
  const gRef = useRef<SVGGElement | null>(null);
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  const [transform, setTransform] = useState({ x: 0, y: 0, scale: 1 });
  const config = useMemo(() => ({ ...defaultConfig, ...userConfig }), [userConfig]);

  // レイアウトサイズの計算
  const layoutSize = useMemo(() => {
    const baseSize = Math.max(dimensions.width, dimensions.height);
    return {
      width: baseSize * 1.5,  // レイアウト領域を広げる
      height: baseSize * 1.5
    };
  }, [dimensions]);

  const processedData = useMemo(() => {
    if (!data.length) return [];

    const maxCount = d3.max(data, d => d.count) || 1;
    const fontSize = d3.scaleSqrt()
      .domain([0, maxCount])
      .range([config.minFontSize, config.maxFontSize]);

    // JSONデータの重要度に基づいてソート
    return data
      .filter(d => d.text && shouldIncludeWord(d.text.trim()))
      .map(d => ({
        text: d.text.trim(),
        size: fontSize(d.count),
        count: d.count,
        importance: calculateWordImportance(d.text.trim(), d.count), // 重要度を計算
        color: config.colors.wordColors[
          Math.floor(Math.random() * config.colors.wordColors.length)
        ],
        rotate: config.layout.rotation.angles[
          Math.floor(Math.random() * config.layout.rotation.angles.length)
        ]
      }))
      .sort((a, b) => b.importance - a.importance) // 重要度でソート
      .slice(0, 100); // 表示する単語数を制限
  }, [data, config]);

  const updateLayout = useCallback(() => {
    if (!gRef.current || !processedData.length || !layoutSize.width) return;

    const layout = cloud()
      .size([layoutSize.width, layoutSize.height])
      .words(processedData)
      .padding(config.layout.padding)
      .rotate(d => d.rotate!)
      .font('Noto Sans JP')
      .fontSize(d => d.size!)
      .spiral(config.layout.spiral)
      .on('end', words => {
        const g = d3.select(gRef.current);
        
        // 既存の要素を更新
        const texts = g.selectAll<SVGTextElement, any>('text')
          .data(words, (d: any) => d.text);

        // 削除
        texts.exit().remove();

        // 新規追加
        const textsEnter = texts.enter()
          .append('text')
          .attr('text-anchor', 'middle')
          .style('font-family', 'Noto Sans JP')
          .style('cursor', 'pointer')
          .style('opacity', 0);

        // 更新（既存 + 新規）
        texts.merge(textsEnter)
          .text(d => d.text)
          .style('font-size', d => `${d.size}px`)
          .style('fill', d => d.color)
          .attr('transform', d => `translate(${d.x},${d.y})rotate(${d.rotate})`)
          .transition()
          .duration(200)
          .style('opacity', 1);

        // イベントハンドラ
        g.selectAll('text')
          .on('mouseover', function() {
            d3.select(this)
              .transition()
              .duration(100)
              .style('fill', config.colors.highlightColor);
          })
          .on('mouseout', function(_, d) {
            d3.select(this)
              .transition()
              .duration(100)
              .style('fill', d.color);
          })
          .on('click', (_, d) => {
            if (onWordClick) onWordClick(d);
          });
      });

    layout.start();
  }, [processedData, layoutSize, config, onWordClick]);

  // コンテナのリサイズ監視
  useEffect(() => {
    const handleResize = debounce(() => {
      if (containerRef.current) {
        const { width, height } = containerRef.current.getBoundingClientRect();
        setDimensions({ width, height });
      }
    }, 100);

    handleResize();
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
      handleResize.cancel();
    };
  }, []);

  // SVG初期化とレイアウト更新
  useEffect(() => {
    if (!svgRef.current || !dimensions.width) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll('*').remove();

    // シンプルな四角形の背景を追加
    svg.append('rect')
      .attr('width', dimensions.width)
      .attr('height', dimensions.height)
      .attr('fill', '#ffffff')
      .attr('rx', 8)  // 軽い角丸
      .attr('ry', 8);

    gRef.current = svg.append('g')
      .attr('transform', `translate(${dimensions.width/2},${dimensions.height/2})`)
      .node();

    updateLayout();
  }, [dimensions, updateLayout]);

  return (
    <div ref={containerRef} className="w-full h-full min-h-[400px] bg-white rounded-lg shadow-lg relative">
      <TransformWrapper
        initialScale={1}
        minScale={0.1}
        maxScale={3}
        centerOnInit={true}
        wheel={{ step: 0.05 }}
        onTransformed={(_, state) => {
          setTransform({
            x: state.positionX,
            y: state.positionY,
            scale: state.scale
          });
        }}
      >
        {({ zoomIn, zoomOut, resetTransform }) => (
          <>
            <div className="absolute top-4 left-4 z-10 flex flex-col gap-2">
              <button
                onClick={() => zoomIn()}
                className="w-8 h-8 bg-white/90 border border-gray-300 rounded shadow hover:bg-gray-50"
              >
                +
              </button>
              <button
                onClick={() => zoomOut()}
                className="w-8 h-8 bg-white/90 border border-gray-300 rounded shadow hover:bg-gray-50"
              >
                -
              </button>
              <button
                onClick={() => resetTransform()}
                className="w-8 h-8 bg-white/90 border border-gray-300 rounded shadow hover:bg-gray-50"
              >
                R
              </button>
            </div>
            <TransformComponent
              wrapperClass="!w-full !h-full"
              contentClass="!w-full !h-full"
            >
              <svg
                ref={svgRef}
                width={dimensions.width}
                height={dimensions.height}
                className="w-full h-full"
                style={{ overflow: 'visible' }}
              />
            </TransformComponent>
          </>
        )}
      </TransformWrapper>
    </div>
  );
};
