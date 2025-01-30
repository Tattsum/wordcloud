export interface WordCloudItem {
    text: string;
    count: number;
    fontSize: number;
    color: string;
  }
  
  export interface WordCloudData {
    items: WordCloudItem[];
    metadata: WordCloudMetadata;
  }
  
  export interface WordCloudMetadata {
    maxCount: number;
    minCount: number;
    totalWords: number;
    generatedAt: string;
  }
  
  export interface WordCloudConfig {
    minFontSize: number;
    maxFontSize: number;
    colors: WordCloudColors;
    layout: WordCloudLayout;
  }
  
  export interface WordCloudColors {
    background: string;
    wordColors: string[];
    highlightColor: string;
  }
  
  export interface WordCloudLayout {
    width: number;
    height: number;
    padding: number;
    spiral: 'archimedean' | 'rectangular';
    rotation: WordCloudRotation;
  }
  
  export interface WordCloudRotation {
    angles: number[];
    random: boolean;
  }
  
  // APIレスポンスの型
  export interface ApiResponse<T> {
    success: boolean;
    data?: T;
    error?: string;
  }
  
  // 解析設定の型
  export interface AnalysisConfig {
    excludeWords: string[];
    minWordLength: number;
    caseSensitive: boolean;
    includeNumbers: boolean;
  }
  
  // フィルタリングオプションの型
  export interface FilterOptions {
    minCount: number;
    maxWords: number;
    excludeUsers: string[];
    dateRange?: {
      start: string;
      end: string;
    };
  }
