type Query = {
  Text: string;
  Column: string;
  ExactMatch: boolean;

  // Services
  Datawells: boolean;
  GithubRecon: boolean;
};

type LeakResult = {
  Duration: number;
  Error: string;
  Rows: Array<Record<string, string>>;
  LimitHit: boolean;
};

type GithubResult = {
  Duration: number;
  Error: string;

  EmailResult: any;
  UsernameResult: any;
};

type Result = {
  Id: string;
  Status: "pending" | "completed";
  Date: string;
  Query: Query;
  LeakResult: LeakResult;
  GithubResult: GithubResult;
};

type HistoryItem = {
  Id: string;
  Status: "pending" | "completed";
  Date: string;
  Query: Query;
  Results: number;
};

type History = HistoryItem[];

type ServerSettings = {
  Folders: string[];
  CacheFolder: string;
  Limit: number;
  MinimumQueryLength: number;

  GithubRecon: boolean;
  GithubTokenLoaded: boolean;
  GithubDeepMode: boolean;
};

type Server = {
  Settings: ServerSettings;

  Dataleaks: Dataleak[];

  TotalRows: number;
  TotalDataleaks: number;
  TotalSize: number;
};

type Dataleak = {
  Name: string;
  Columns: string[];
  Length: number;
  Size: number;
};

export type {
  Query,
  LeakResult,
  History,
  HistoryItem,
  GithubResult,
  Result,
  ServerSettings,
  Server,
  Dataleak,
};
