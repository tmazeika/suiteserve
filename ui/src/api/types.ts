export type Id = string | number;

export interface Entity {
  id: Id;
}

export interface VersionedEntity {
  version: number;
}

export interface SoftDeleteEntity {
  deleted?: boolean;
  deleted_at?: number;
}

export interface Attachment extends Entity, VersionedEntity, SoftDeleteEntity {
  suite_id?: Id;
  case_id?: Id;
  filename: string;
  content_type: string;
  size: number;
  timestamp: number;
}

export enum SuiteStatus {
  STARTED = 'started',
  FINISHED = 'finished',
  DISCONNECTED = 'disconnected',
}

export enum SuiteResult {
  PASSED = 'passed',
  FAILED = 'failed',
}

export interface Suite extends Entity, VersionedEntity, SoftDeleteEntity {
  name?: string;
  tags?: string[];
  planned_cases?: number;
  status: SuiteStatus | string;
  result?: SuiteResult | string;
  disconnected_at?: number;
  started_at: number;
  finished_at?: number;
}

export interface SuitePage {
  more: boolean;
  suites: Suite[];
}

export enum CaseStatus {
  CREATED = 'created',
  STARTED = 'started',
  FINISHED = 'finished',
}

export enum CaseResult {
  PASSED = 'passed',
  FAILED = 'failed',
  SKIPPED = 'skipped',
  ABORTED = 'aborted',
  ERRORED = 'errored',
}

type JsonValue =
  | string
  | number
  | boolean
  | null
  | Map<string, JsonValue>
  | Array<JsonValue>;

export interface Case extends Entity, VersionedEntity {
  suite_id: Id;
  name?: string;
  description?: string;
  tags?: string[];
  idx: number;
  args?: {
    [key: string]: JsonValue;
  };
  status: CaseStatus | string;
  result?: CaseResult | string;
  created_at: number;
  started_at?: number;
  finished_at?: number;
}

export enum LogLevelType {
  TRACE = 'trace',
  DEBUG = 'debug',
  INFO = 'info',
  WARN = 'warn',
  ERROR = 'error',
}

export interface LogLine extends Entity {
  case_id: Id;
  idx: number;
  level: LogLevelType | string;
  trace?: string;
  message?: string;
  timestamp: number;
}

async function getJson<T>(url: string): Promise<T> {
  const resp = await fetch(url);
  if (!resp.ok) {
    throw new Error(resp.statusText);
  }
  return resp.json();
}

export async function getAttachment(id: string): Promise<Attachment> {
  return getJson(`/v1/attachments/${encodeURIComponent(id)}`);
}

export async function getSuiteAttachments(id: string): Promise<Attachment[]> {
  return getJson(`/v1/attachments?suite=${encodeURIComponent(id)}`);
}

export async function getCaseAttachments(id: string): Promise<Attachment[]> {
  return getJson(`/v1/attachments?case=${encodeURIComponent(id)}`);
}

export async function getAllAttachments(): Promise<Attachment[]> {
  return getJson('/v1/attachments');
}

export async function getSuite(id: string): Promise<Suite> {
  return getJson(`/v1/suites/${encodeURIComponent(id)}`);
}

export function watchSuites(): EventSource {
  return new EventSource('/v1/suites?watch=true');
}

export async function getSuitePage(): Promise<SuitePage> {
  return getJson('/v1/suites');
}

export async function getSuitePageAfter(id: string): Promise<SuitePage> {
  return getJson(`/v1/suites?after=${encodeURIComponent(id)}`);
}

export async function getCase(id: string): Promise<Case> {
  return getJson(`/v1/cases/${encodeURIComponent(id)}`);
}

export async function getLogLine(id: string): Promise<LogLine> {
  return getJson(`/v1/logs/${encodeURIComponent(id)}`);
}
