import Papa from 'papaparse';
import { SlackMessageCSV } from '@/lib/types/slack';

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
