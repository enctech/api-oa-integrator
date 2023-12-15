import { useSession } from "../context/session-context";
import React, { ReactNode } from "react"; // component that take children and return children if user is admin

export const AdminOnly: React.FC<{ children: ReactNode }> = ({ children }) => {
  const { session } = useSession();
  if (session?.permissions?.includes("admin")) {
    return <>{children}</>;
  }
  return <></>;
};
