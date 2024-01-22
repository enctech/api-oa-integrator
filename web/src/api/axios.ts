import axios from "axios";
import authMiddleware from "../middlewares/auth-middlewares";

const instance = axios.create({
  baseURL: window.location.origin.replace(":3000", ":1323"),
  timeout: 10000,
});

export const internal = axios.create({
  baseURL: window.location.origin.replace(":3000", ":1323"),
  timeout: 10000,
});

instance.interceptors.request.use(authMiddleware);

export default instance;
