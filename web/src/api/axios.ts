import axios from "axios";
import authMiddleware from "../middlewares/auth-middlewares";

const instance = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  timeout: 10000,
});

export const internal = axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  timeout: 10000,
});

instance.interceptors.request.use(authMiddleware);

export default instance;
