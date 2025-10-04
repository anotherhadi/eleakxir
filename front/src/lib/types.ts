export type Query = {
  Text: string;
  Column: string;
  ExactMatch: boolean;

  // Services
  Datawells: boolean;
  GithubRecon: boolean;
  GravatarRecon: boolean;
};

export type LeakResult = {
  Duration: number;
  Error: string;
  Rows: Array<Record<string, string>>;
  LimitHit: boolean;
  Inactive: boolean;
};

export type GithubResult = {
  Duration: number;
  Error: string;
  Inactive: boolean;

  EmailResult: any;
  UsernameResult: any;
};

export type GravatarResult = {
  Duration: number;
  Error: string;
  Inactive: boolean;

  Results: any;
};

export type Result = {
  Id: string;
  Status: "queued" | "pending" | "completed" | "error";
  Date: string;
  Query: Query;
  ResultsCount: number;

  LeakResult: LeakResult;
  GithubResult: GithubResult;
  GravatarResult: GravatarResult;
};

export type History = Result[];

export type ServerSettings = {
  Folders: string[];
  CacheFolder: string;
  Limit: number;
  MinimumQueryLength: number;

  GithubRecon: boolean;
  GithubTokenLoaded: boolean;
  GithubDeepMode: boolean;
  GravatarRecon: boolean;
};

export type Server = {
  Settings: ServerSettings;

  Dataleaks: Dataleak[];

  TotalRows: number;
  TotalDataleaks: number;
  TotalSize: number;
};

export type Dataleak = {
  Name: string;
  Columns: string[];
  Length: number;
  Size: number;
  ModTime: string;
  Path: string;
};
