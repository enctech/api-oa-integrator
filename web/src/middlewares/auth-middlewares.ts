import { InternalAxiosRequestConfig } from "axios";
import { resolvePromise } from "../utils/promise";
import { internal } from "../api/axios";

const authMiddleware = async (config: InternalAxiosRequestConfig<any>) => {
  const response = await resolvePromise(
    internal.post("/auth/login", {
      username: "string",
      password: "string",
    }),
  );
  console.log("ERROR", response?.error);
  const token = response?.result?.data.token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
};

export default authMiddleware;
