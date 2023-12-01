import { InternalAxiosRequestConfig } from "axios";

const authMiddleware = async (config: InternalAxiosRequestConfig<any>) => {
  const storedUserData = sessionStorage.getItem("userData");
  if (!storedUserData) return config;
  const token = JSON.parse(storedUserData)?.token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
};

export default authMiddleware;
