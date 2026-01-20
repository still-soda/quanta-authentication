const BASE_URL = 'http://localhost:8080';
const r = (uri: string, option: RequestInit) =>
   fetch(BASE_URL + uri, {
      ...option,
      headers: {
         'Content-Type': 'application/json',
         ...option?.headers,
      },
   }).then((res) => res.json());
const b = (token: string) => ({
   Authorization: `Bearer ${token}`,
});

const secret = 'supersecret';

async function loginToAdmin() {
   const res = await r('/auth/login', {
      method: 'POST',
      body: JSON.stringify({
         student_id: '20231003059',
         password: '123456',
      }),
   });

   console.log(res.code.toString().startsWith('2') ? 'PASS' : 'FAILED');

   return {
      freshToken: res.data.refresh_token,
      accessToken: res.data.access_token,
      userId: res.data.user.id,
   };
}

async function createClient(token: string, id: string) {
   const res = await r('/oauth/clients', {
      method: 'POST',
      headers: b(token),
      body: JSON.stringify({
         name: 'Test Client',
         domain: 'http://localhost',
         user_id: id,
         is_public: false,
         secret: secret,
      }),
   });

   console.log(res.code.toString().startsWith('2') ? 'PASS' : 'FAILED');

   return res.data.client_id;
}

async function removeClient(token: string, clientId: string) {
   const res = await r(`/oauth/clients/${clientId}`, {
      method: 'DELETE',
      headers: b(token),
   });

   console.log(res.code.toString().startsWith('2') ? 'PASS' : 'FAILED');
}

async function authorizeClient(clientId: string, token: string) {
   const params = new URLSearchParams({
      client_id: clientId,
      response_type: 'code',
      redirect_uri: 'http://localhost/callback',
      scope: 'openid',
      state: 'xyz',
   });
   const res = await fetch(`${BASE_URL}/oauth/authorize?${params.toString()}`, {
      credentials: 'include',
      redirect: 'manual',
      headers: {
         Cookie: `access_token=${token}`,
      },
   });

   const location = res.headers.get('Location');
   // console.log(`Redirect to ${location}`)
   const code = location ? new URL(location).searchParams.get('code') : null;
   console.log(code ? 'PASS' : 'FAILED');

   return code ?? '';
}

async function tokenFromCode(clientId: string, code: string) {
   const params = new URLSearchParams({
      grant_type: 'authorization_code',
      code: code,
      redirect_uri: 'http://localhost/callback',
      client_id: clientId,
      client_secret: secret,
   });
   const res = await r('/oauth/token', {
      method: 'POST',
      body: params,
      headers: {
         'Content-Type': 'application/x-www-form-urlencoded',
      }
   });

   // console.log(res)
   console.log('access_token' in res ? 'PASS' : 'FAILED');

   return {
      accessToken: res.access_token,
      idToken: res.id_token,
      refreshToken: res.refresh_token,
   };
}

async function verifyIdToken(token: string) {
   const [headerB64, payload, signature] = token.split('.');
   const decodedHeader = JSON.parse(atob(headerB64));

   const res1 = await r('/.well-known/jwks.json', { method: 'GET' });
   
   // 找到匹配的 key
   const jwk = res1.keys.find((k: any) => k.kid === decodedHeader.kid);
   const { n, e } = jwk;

   const cryptoKey = await crypto.subtle.importKey(
      'jwk',
      { kty: 'RSA', n: n, e: e, alg: 'RS256', ext: true },
      { name: 'RSASSA-PKCS1-v1_5', hash: 'SHA-256' },
      false,
      ['verify']
   );
   
   const encoder = new TextEncoder();
   const data = encoder.encode([headerB64, payload].join('.'));
   
   const b64 = atob(signature.replace(/-/g, '+').replace(/_/g, '/'))
   const signatureArray = Uint8Array.from(b64, c => c.charCodeAt(0));
   const isValid = await crypto.subtle.verify(
      'RSASSA-PKCS1-v1_5',
      cryptoKey,
      signatureArray,
      data
   );

   console.log(isValid ? 'PASS' : 'FAILED');
}

async function getUserInfo(token: string) {
   const res = await r('/oauth/userinfo', {
      method: 'GET',
      headers: b(token),
   });

   console.log(res.code.toString().startsWith('2') ? 'PASS' : 'FAILED');
}

async function logout(token: string) {
   const res = await r('/oauth/logout', {
      method: 'POST',
      headers: b(token),
   });

   console.log(res.code.toString().startsWith('2') ? 'PASS' : 'FAILED');
}

async function main() {
   // 创建阶段
   const cred = await loginToAdmin();
   const clientId = await createClient(cred.accessToken, cred.userId);

   // 授权阶段
   const code = await authorizeClient(clientId, cred.accessToken);
   const token = await tokenFromCode(clientId, code);
   await verifyIdToken(token.idToken);

   // 获取用户信息
   await getUserInfo(token.accessToken);

   // 清理阶段
   await logout(cred.accessToken);

   await removeClient(cred.accessToken, clientId);
}

main();
