'use client';

import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { AuthUser, getCurrentUser, signIn, signOut as cognitoSignOut, SignInParams } from '@/lib/cognito';

interface AuthContextType {
  user: AuthUser | null;
  loading: boolean;
  signIn: (params: SignInParams) => Promise<void>;
  signOut: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [loading, setLoading] = useState(true);
  
  // 認証が無効の場合はスキップ
  const authEnabled = process.env.NEXT_PUBLIC_AUTH_ENABLED === 'true';

  const refreshUser = async () => {
    if (!authEnabled) {
      setUser(null);
      return;
    }
    
    try {
      const currentUser = await getCurrentUser();
      setUser(currentUser);
    } catch (error) {
      setUser(null);
    }
  };

  useEffect(() => {
    refreshUser().finally(() => setLoading(false));
  }, []);

  const handleSignIn = async (params: SignInParams) => {
    if (!authEnabled) {
      throw new Error('認証機能は現在無効です。.env.localでNEXT_PUBLIC_AUTH_ENABLED=trueを設定してください。');
    }
    
    try {
      await signIn(params);
      await refreshUser();
    } catch (error) {
      console.error('Sign in error:', error);
      throw error;
    }
  };

  const handleSignOut = () => {
    if (!authEnabled) {
      return;
    }
    cognitoSignOut();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        signIn: handleSignIn,
        signOut: handleSignOut,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};
