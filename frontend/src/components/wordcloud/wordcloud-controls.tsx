import React from 'react';
import { WordCloudConfig, WordCloudItem } from '@/lib/types/wordcloud';
import { parseCSV, parseJSON } from '@/lib/utils/file-parser';

interface WordCloudControlsProps {
  config: WordCloudConfig;
  onConfigChange: (config: WordCloudConfig) => void;
  onDataUpload: (data: WordCloudItem[]) => void;
}

export const WordCloudControls: React.FC<WordCloudControlsProps> = ({
  config,
  onConfigChange,
  onDataUpload
}) => {
  const handleConfigChange = (
    key: keyof WordCloudConfig,
    value: any
  ) => {
    onConfigChange({
      ...config,
      [key]: value
    });
  };

  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    try {
      let data: WordCloudItem[];
      if (file.name.endsWith('.json')) {
        data = await parseJSON(file);
      } else if (file.name.endsWith('.csv')) {
        const csvData = await parseCSV(file);
        // CSVデータをWordCloudItem形式に変換
        data = csvData.map(item => ({
          text: item.text || '',
          count: typeof item.count === 'number' ? item.count : 1,
          fontSize: 0, // これは後で計算される
          color: '#000096' // デフォルトカラー
        }));
      } else {
        throw new Error('サポートされていないファイル形式です');
      }
      onDataUpload(data);
    } catch (error) {
      console.error('ファイルの解析に失敗しました:', error);
      alert('ファイルの解析に失敗しました');
    }
  };

  const handleDataUpload = (newData: WordCloudItem[]) => {
    onDataUpload(newData);
  };

  return (
    <div className="space-y-4 p-4 bg-white rounded-lg shadow-lg max-h-[500px] overflow-y-auto">
      <div>
        <label className="block text-sm font-medium text-gray-700">
          データファイルのアップロード
        </label>
        <input
          type="file"
          accept=".csv,.json"
          onChange={handleFileUpload}
          className="mt-1 block w-full text-sm text-gray-500
            file:mr-4 file:py-2 file:px-4
            file:rounded-full file:border-0
            file:text-sm file:font-semibold
            file:bg-blue-50 file:text-blue-700
            hover:file:bg-blue-100"
        />
        <p className="mt-1 text-sm text-gray-500">
          CSVまたはJSONファイルをアップロードしてください
        </p>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700">
          最小フォントサイズ
        </label>
        <input
          type="range"
          min="10"
          max="40"
          value={config.minFontSize}
          onChange={(e) => handleConfigChange('minFontSize', parseInt(e.target.value))}
          className="w-full"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700">
          最大フォントサイズ
        </label>
        <input
          type="range"
          min="40"
          max="120"
          value={config.maxFontSize}
          onChange={(e) => handleConfigChange('maxFontSize', parseInt(e.target.value))}
          className="w-full"
        />
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700">
          回転
        </label>
        <select
          value={config.layout.rotation.random ? 'random' : 'fixed'}
          onChange={(e) => handleConfigChange('layout', {
            ...config.layout,
            rotation: {
              angles: e.target.value === 'random' ? [-90, -45, 0, 45, 90] : [0],
              random: e.target.value === 'random'
            }
          })}
          className="w-full mt-1 rounded-md border-gray-300 shadow-sm focus:border-primary focus:ring-primary"
        >
          <option value="fixed">固定</option>
          <option value="random">ランダム</option>
        </select>
      </div>
    </div>
  );
};
