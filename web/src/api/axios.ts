import axios from "axios";
import authMiddleware from "../middlewares/auth-middlewares";

const baseURL = "http://localhost:1323";

const instance = axios.create({
  baseURL: baseURL, // Replace with your API URL
  timeout: 10000, // Request timeout in milliseconds
});

export const internal = axios.create({
  baseURL: baseURL, // Replace with your API URL
  timeout: 10000, // Request timeout in milliseconds
});

instance.interceptors.request.use(authMiddleware);

export default instance;
