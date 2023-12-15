// SessionContext.tsx
import { createContext, useContext, useEffect, useState } from "react";

interface User {
  username: string;
  userId: string;
  token: string;
  refreshToken: string;
  permissions: string[];
}

interface SessionContextProps {
  session: Partial<User> | null;
  login: (user: User) => void;
  logout: () => void;
}

const SessionContext = createContext<SessionContextProps | undefined>(
  undefined,
);

export const useSession = (): SessionContextProps => {
  const context = useContext(SessionContext);
  if (!context) {
    throw new Error("useSession must be used within a SessionProvider");
  }
  return context;
};

export const SessionProvider = ({ children }: any) => {
  const [session, setSession] = useState<User | null>(null);

  const login = (user: User) => {
    // Simulating a login process, set the session data
    setSession(user);
  };

  const logout = () => {
    // Simulating a logout process, clear the session data
    sessionStorage.removeItem("userData");
    setSession(null);
  };

  useEffect(() => {
    const storedUserData = sessionStorage.getItem("userData");
    if (!storedUserData) return;
    setSession(JSON.parse(storedUserData));
  }, []);

  const contextValue: SessionContextProps = {
    session,
    login,
    logout,
  };

  return (
    <SessionContext.Provider value={contextValue}>
      {children}
    </SessionContext.Provider>
  );
};
