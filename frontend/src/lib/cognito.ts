import { CognitoUserPool, CognitoUser, AuthenticationDetails, CognitoUserSession, CognitoUserAttribute } from 'amazon-cognito-identity-js';

// Cognito設定
const poolData = {
  UserPoolId: process.env.NEXT_PUBLIC_COGNITO_USER_POOL_ID || '',
  ClientId: process.env.NEXT_PUBLIC_COGNITO_CLIENT_ID || '',
};

let userPool: CognitoUserPool | null = null;

// UserPoolの初期化
const getUserPool = () => {
  if (!userPool && poolData.UserPoolId && poolData.ClientId) {
    userPool = new CognitoUserPool(poolData);
  }
  return userPool;
};

export interface AuthUser {
  username: string;
  email: string;
  sub: string;
}

export interface SignInParams {
  email: string;
  password: string;
}

export interface SignUpParams {
  email: string;
  password: string;
  name?: string;
}

/**
 * サインイン
 */
export const signIn = (params: SignInParams): Promise<CognitoUserSession> => {
  const { email, password } = params;

  return new Promise((resolve, reject) => {
    const pool = getUserPool();
    if (!pool) {
      reject(new Error('Cognito User Pool not configured'));
      return;
    }

    const authenticationDetails = new AuthenticationDetails({
      Username: email,
      Password: password,
    });

    const cognitoUser = new CognitoUser({
      Username: email,
      Pool: pool,
    });

    cognitoUser.authenticateUser(authenticationDetails, {
      onSuccess: (session) => {
        resolve(session);
      },
      onFailure: (err) => {
        reject(err);
      },
    });
  });
};

/**
 * サインアウト
 */
export const signOut = (): void => {
  const pool = getUserPool();
  if (!pool) return;

  const cognitoUser = pool.getCurrentUser();
  if (cognitoUser) {
    cognitoUser.signOut();
  }
  
  // ローカルストレージからトークンを削除
  if (typeof window !== 'undefined') {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('id_token');
    localStorage.removeItem('refresh_token');
  }
};

/**
 * 現在のユーザーセッションを取得
 */
export const getCurrentSession = (): Promise<CognitoUserSession> => {
  return new Promise((resolve, reject) => {
    const pool = getUserPool();
    if (!pool) {
      reject(new Error('Cognito User Pool not configured'));
      return;
    }

    const cognitoUser = pool.getCurrentUser();
    if (!cognitoUser) {
      reject(new Error('No current user'));
      return;
    }

    cognitoUser.getSession((err: Error | null, session: CognitoUserSession | null) => {
      if (err || !session) {
        reject(err || new Error('No session'));
        return;
      }

      if (!session.isValid()) {
        reject(new Error('Session is not valid'));
        return;
      }

      resolve(session);
    });
  });
};

/**
 * 現在のユーザー情報を取得
 */
export const getCurrentUser = async (): Promise<AuthUser | null> => {
  try {
    const session = await getCurrentSession();
    const idToken = session.getIdToken();
    const payload = idToken.payload;

    return {
      username: payload['cognito:username'] || payload.email,
      email: payload.email,
      sub: payload.sub,
    };
  } catch (error) {
    return null;
  }
};

/**
 * IDトークンを取得
 */
export const getIdToken = async (): Promise<string | null> => {
  try {
    const session = await getCurrentSession();
    return session.getIdToken().getJwtToken();
  } catch (error) {
    return null;
  }
};

/**
 * アクセストークンを取得
 */
export const getAccessToken = async (): Promise<string | null> => {
  try {
    const session = await getCurrentSession();
    return session.getAccessToken().getJwtToken();
  } catch (error) {
    return null;
  }
};

/**
 * サインアップ
 */
export const signUp = (params: SignUpParams): Promise<any> => {
  const { email, password, name } = params;

  return new Promise((resolve, reject) => {
    const pool = getUserPool();
    if (!pool) {
      reject(new Error('Cognito User Pool not configured'));
      return;
    }

    const attributeList = [];

    if (name) {
      attributeList.push(
        new CognitoUserAttribute({
          Name: 'name',
          Value: name,
        })
      );
    }

    pool.signUp(email, password, attributeList, [], (err, result) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(result);
    });
  });
};

/**
 * 確認コードを送信してアカウントを確認
 */
export const confirmSignUp = (email: string, code: string): Promise<any> => {
  return new Promise((resolve, reject) => {
    const pool = getUserPool();
    if (!pool) {
      reject(new Error('Cognito User Pool not configured'));
      return;
    }

    const cognitoUser = new CognitoUser({
      Username: email,
      Pool: pool,
    });

    cognitoUser.confirmRegistration(code, true, (err, result) => {
      if (err) {
        reject(err);
        return;
      }
      resolve(result);
    });
  });
};

/**
 * パスワードリセットのリクエスト
 */
export const forgotPassword = (email: string): Promise<any> => {
  return new Promise((resolve, reject) => {
    const pool = getUserPool();
    if (!pool) {
      reject(new Error('Cognito User Pool not configured'));
      return;
    }

    const cognitoUser = new CognitoUser({
      Username: email,
      Pool: pool,
    });

    cognitoUser.forgotPassword({
      onSuccess: (data) => {
        resolve(data);
      },
      onFailure: (err) => {
        reject(err);
      },
    });
  });
};

/**
 * パスワードの確認とリセット
 */
export const confirmPassword = (
  email: string,
  code: string,
  newPassword: string
): Promise<any> => {
  return new Promise((resolve, reject) => {
    const pool = getUserPool();
    if (!pool) {
      reject(new Error('Cognito User Pool not configured'));
      return;
    }

    const cognitoUser = new CognitoUser({
      Username: email,
      Pool: pool,
    });

    cognitoUser.confirmPassword(code, newPassword, {
      onSuccess: () => {
        resolve('Password confirmed!');
      },
      onFailure: (err) => {
        reject(err);
      },
    });
  });
};
