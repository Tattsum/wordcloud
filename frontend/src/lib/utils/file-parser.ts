import Papa from 'papaparse';
import { SlackMessageCSV } from '@/lib/types/slack';
import { WordCloudItem } from '@/lib/types/wordcloud';

export const parseCSV = (file: File): Promise<SlackMessageCSV[]> => {
  return new Promise((resolve, reject) => {
    Papa.parse(file, {
      header: true,
      skipEmptyLines: true,
      dynamicTyping: true,
      complete: (results) => {
        if (results.errors.length > 0) {
          reject(new Error('CSVファイルの解析中にエラーが発生しました'));
          return;
        }
        resolve(results.data as SlackMessageCSV[]);
      },
      error: (error) => {
        reject(new Error(`CSVファイルの解析に失敗しました: ${error.message}`));
      },
    });
  });
};

export const parseJSON = (file: File): Promise<WordCloudItem[]> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (event) => {
      try {
        const json = JSON.parse(event.target?.result as string);
        resolve(json as WordCloudItem[]);
      } catch (error) {
        reject(new Error('JSONファイルの解析に失敗しました'));
      }
    };
    reader.onerror = () => {
      reject(new Error('JSONファイルの読み込みに失敗しました'));
    };
    reader.readAsText(file);
  });
}; 
