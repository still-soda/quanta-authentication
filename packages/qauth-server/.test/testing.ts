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
   const code = location ? new URL(location).searchParams.get('code') : null;
   console.log(code ? 'PASS' : 'FAILED');

   return code;
}

async function main() {
   const cred = await loginToAdmin();
   const clientId = await createClient(cred.accessToken, cred.userId);

   await authorizeClient(clientId, cred.accessToken);

   await removeClient(cred.accessToken, clientId);
}

main();
