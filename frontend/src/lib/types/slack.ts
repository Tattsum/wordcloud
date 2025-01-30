export interface SlackMessage {
    timestamp: string;
    userId: string;
    username: string;
    message: string;
    threadTs: string | null;
  }
  
  export interface SlackMessageCSV {
    Timestamp: string;
    UserID: string;
    Username: string;
    Message: string;
    ThreadTS: string;
  }
  
  export interface SlackChannel {
    id: string;
    name: string;
  }
  
  // ファイルアップロード関連の型
  export interface FileUploadState {
    isUploading: boolean;
    progress: number;
    error: string | null;
  }
