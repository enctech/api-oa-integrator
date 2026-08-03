import axios from "axios";
import authMiddleware from "../middlewares/auth-middlewares";

const instance = axios.create({
  baseURL: window.location.origin + "/api/",
  timeout: 10000,
});

instance.interceptors.request.use(authMiddleware);

export default instance;
