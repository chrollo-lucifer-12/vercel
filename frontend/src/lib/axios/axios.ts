import axios from "axios";
import { clientEnv } from "../env/client";

export const axiosInstance = axios.create({
  baseURL: clientEnv.NEXT_PUBLIC_BACKEND_ENDPOINT,
  withCredentials: true,
});
