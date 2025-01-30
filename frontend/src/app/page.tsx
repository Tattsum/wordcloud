'use client';

import React, { useState } from 'react';
import { FileUpload } from '@/components/ui/file-upload';
import { WordCloudCanvas } from '@/components/wordcloud/wordcloud-canvas';
import { SlackMessageCSV, WordCloudItem } from '@/lib/types/wordcloud';

export default function Home() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [wordCloudData, setWordCloudData] = useState<WordCloudItem[]>([]);
  const [wordCountMap, setWordCountMap] = useState<Map<string, number>>(new Map());
  const [styleMap, setStyleMap] = useState<Map<string, { color: string }>>(new Map());

  const updateWordCloud = (wordCounts: Map<string, number>, styles: Map<string, { color: string }>) => {
    const cloudData: WordCloudItem[] = Array.from(wordCounts.entries())
      .map(([text, count]) => {
        const style = styles.get(text) || { color: '#2563eb' };
        // カウントに基づいてフォントサイズを動的に計算
        const fontSize = calculateFontSize(count);
        return {
          text,
          count,
          fontSize,
          color: style.color
        };
      })
      .filter(item => item.text.length > 1);

    setWordCloudData(cloudData);
  };

  const calculateFontSize = (count: number): number => {
    // カウントに基づいて適切なフォントサイズを計算
    const minCount = 1;
    const maxCount = Math.max(...Array.from(wordCountMap.values()));
    const minSize = 14;
    const maxSize = 80;
    
    if (maxCount === minCount) return minSize;
    
    const size = minSize + (count - minCount) * (maxSize - minSize) / (maxCount - minCount);
    return Math.round(size);
  };

  const handleFileProcess = async (
    data: SlackMessageCSV[] | WordCloudItem[],
    fileType: 'csv' | 'json'
  ) => {
    try {
      setLoading(true);
      setError(null);

      const currentWordCounts = new Map(wordCountMap);
      const currentStyles = new Map(styleMap);

      if (fileType === 'json') {
        // JSONデータの場合、スタイル情報を更新
        const jsonData = data as WordCloudItem[];
        jsonData.forEach(item => {
          currentStyles.set(item.text, { color: item.color });
          // 既存のカウントがあれば加算、なければ新規設定
          const currentCount = currentWordCounts.get(item.text) || 0;
          currentWordCounts.set(item.text, currentCount + item.count);
        });
      } else {
        // CSVデータの場合、単語カウントを更新
        const csvData = data as SlackMessageCSV[];
        csvData.forEach(row => {
          const words = row.Message.split(/\s+/);
          words.forEach(word => {
            const currentCount = currentWordCounts.get(word) || 0;
            currentWordCounts.set(word, currentCount + 1);
          });
        });
      }

      setWordCountMap(currentWordCounts);
      setStyleMap(currentStyles);
      updateWordCloud(currentWordCounts, currentStyles);

    } catch (err) {
      setError('データの処理中にエラーが発生しました');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleClearData = () => {
    setWordCloudData([]);
    setWordCountMap(new Map());
    setStyleMap(new Map());
    setError(null);
  };

  return (
    <main className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-8">
          Slack ワードクラウド生成
        </h1>

        <div className="space-y-8">
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold text-gray-800">
                ファイルをアップロード
              </h2>
              {wordCloudData.length > 0 && (
                <button
                  onClick={handleClearData}
                  className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600 transition-colors"
                >
                  データをクリア
                </button>
              )}
            </div>
            <FileUpload
              onFileSelect={() => setError(null)}
              onFileProcess={handleFileProcess}
            />
            <p className="mt-2 text-sm text-gray-500">
              CSVファイル: テキストデータとして処理されます<br />
              JSONファイル: スタイル情報として処理されます
            </p>
          </div>

          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-4">
              <p className="text-red-600">{error}</p>
            </div>
          )}

          {loading && (
            <div className="flex justify-center items-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          )}

          {wordCloudData.length > 0 && (
            <div className="bg-white rounded-lg shadow">
              <div className="p-6">
                <h2 className="text-xl font-semibold text-gray-800 mb-4">
                  ワードクラウド
                </h2>
                <div className="w-full h-[600px]">
                  <WordCloudCanvas
                    data={wordCloudData}
                    onWordClick={(word) => {
                      console.log('Selected word:', word);
                    }}
                  />
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </main>
  );
}
