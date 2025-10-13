import { createRequest } from '../index.ts';

const baseUrl = '/im_connections/';
const request = createRequest(baseUrl);

export function getConnectionList() {
  return request<DiceConnection[]>('get', 'list');
}

export function getLagrangeSignInfo() {
  return request<{ result: true; data: SignInfo[] } | { result: false; err: string }>(
    'get',
    'get_lgr_signinfo',
  );
}

export function postAddDiscord(
  token: string,
  proxyURL: string,
  reverseProxyUrl: string,
  reverseProxyCDNUrl: string,
) {
  return request<DiceConnection>(
    'post',
    'addDiscord',
    { token, proxyURL, reverseProxyUrl, reverseProxyCDNUrl },
    'json',
    { timeout: 65000 },
  );
}

export function postAddKook(token: string) {
  return request<DiceConnection>('post', 'addKook', { token }, 'json', {
    timeout: 65000,
  });
}

export function postAddTelegram(token: string, proxyURL: string) {
  return request<DiceConnection>('post', 'addTelegram', { token, proxyURL }, 'json', {
    timeout: 65000,
  });
}

export function postAddGocqSeparate(
  relWorkDir: string,
  connectUrl: string,
  accessToken: string,
  account: string,
) {
  return request<DiceConnection>(
    'post',
    'addGocqSeparate',
    { relWorkDir, connectUrl, accessToken, account },
    'json',
    {
      timeout: 65000,
    },
  );
}

export function postAddOnebot11ReverseWs(account: string, reverseAddr?: string) {
  return request<DiceConnection>('post', 'addOnebot11ReverseWs', { account, reverseAddr }, 'json', {
    timeout: 65000,
  });
}

export function postAddLagrange(
  account: string,
  signServerName: string,
  signServerVersion: string,
  isGocq: boolean,
) {
  return request<DiceConnection>(
    'post',
    'addLagrange',
    { account, signServerName, signServerVersion, isGocq },
    'json',
    {
      timeout: 65000,
    },
  );
}

export function postAddMilky(token: string, wsGateway: string, restGateway: string) {
  return request<DiceConnection>('post', 'addMilky', { token, wsGateway, restGateway }, 'json', {
    timeout: 65000,
  });
}

export function postConnectionDel(id: string) {
  return request<DiceConnection>('post', 'del', { id });
}

export function postConnectionQrcode(id: string) {
  return request<{ img: string }>('post', 'qrcode', { id });
}

export function postConnectSetEnable(id: string, enable: boolean) {
  return request<DiceConnection>('post', 'set_enable', { id, enable });
}

export interface DiceConnection {
  id: string;
  state: number;
  platform: string;
  workDir: string;
  enable: boolean;
  protocolType: string;
  nickname: string;
  userId: number;
  groupNum: number;
  cmdExecutedNum: number;
  cmdExecutedLastTime: number;
  isPublic: boolean;

  adapter: AdapterQQ;
}

interface AdapterQQ {
  DiceServing: boolean;
  connectUrl: string;
  curLoginFailedReason: string;
  curLoginIndex: number;
  loginState: goCqHttpStateCode;
  inPackGoCqHttpLastRestricted: number;
  inPackGoCqHttpProtocol: number;
  inPackGoCqHttpAppVersion: string;
  implementation: string;
  useInPackGoCqhttp: boolean;
  goCqHttpLoginVerifyCode: string;
  goCqHttpLoginDeviceLockUrl: string;
  ignoreFriendRequest: boolean;
  goCqHttpSmsNumberTip: string;
  useSignServer: boolean;
  redVersion: string;
  host: string;
  port: number;
  appID: number;
  isReverse: boolean;
  reverseAddr: string;
  builtinMode: 'gocq' | 'lagrange' | 'lagrange-gocq';
  signServerVer: string;
  signServerName: string;
}

enum goCqHttpStateCode {
  Init = 0,
  InLogin = 1,
  InLoginQrCode = 2,
  LoginSuccessed = 10,
  LoginFailed = 11,
  Closed = 20,
}

interface SignServer {
  name: string;
  url: string;
  latency: number;
  selected: boolean;
  ignored: boolean;
  note: string;
}
export interface SignInfo {
  version: string;
  appinfo: [];
  servers: SignServer[];
  selected: boolean;
  ignored: boolean;
  note: string;
}
