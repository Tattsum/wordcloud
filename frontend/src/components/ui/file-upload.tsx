import React, { useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import { parseCSV, parseJSON } from '@/lib/utils/file-parser';
import { SlackMessageCSV, WordCloudItem } from '@/lib/types/wordcloud';

interface FileUploadProps {
  onFileSelect: () => void;
  onFileProcess: (data: SlackMessageCSV[] | WordCloudItem[], fileType: 'csv' | 'json') => void;
}

export const FileUpload: React.FC<FileUploadProps> = ({
  onFileSelect,
  onFileProcess,
}) => {
  const onDrop = useCallback(
    async (acceptedFiles: File[]) => {
      onFileSelect();

      for (const file of acceptedFiles) {
        try {
          if (file.name.endsWith('.json')) {
            const data = await parseJSON(file);
            onFileProcess(data, 'json');
          } else if (file.name.endsWith('.csv')) {
            const data = await parseCSV(file);
            onFileProcess(data, 'csv');
          } else {
            throw new Error('サポートされていないファイル形式です');
          }
        } catch (error) {
          console.error(`${file.name}の処理中にエラーが発生しました:`, error);
          throw error;
        }
      }
    },
    [onFileSelect, onFileProcess]
  );

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'text/csv': ['.csv'],
      'application/json': ['.json'],
    },
    maxSize: 5 * 1024 * 1024, // 5MB
    multiple: true, // 複数ファイルを許可
  });

  return (
    <div
      {...getRootProps()}
      className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors
        ${
          isDragActive
            ? 'border-primary bg-primary/5'
            : 'border-gray-300 hover:border-primary/50'
        }`}
    >
      <input {...getInputProps()} />
      <div className="flex flex-col items-center space-y-2">
        <svg
          className="w-12 h-12 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
          />
        </svg>
        <div className="text-gray-600">
          {isDragActive ? (
            <p>ここにファイルをドロップ</p>
          ) : (
            <>
              <p>ここにファイルをドラッグ＆ドロップ</p>
              <p className="text-sm text-gray-500">
                または、クリックしてファイルを選択
              </p>
            </>
          )}
        </div>
        <p className="text-sm text-gray-500">
          CSVまたはJSONファイル（最大 5MB）
        </p>
      </div>
    </div>
  );
};
