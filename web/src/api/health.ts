import axios from "./axios";

interface HealthArg {
  facility: string;
  device: string;
}

const sample = {
  db: "up",
  oa: "down",
};

export const getOAHealth = async (arg: HealthArg) => {
  return axios
    .get(`/health?facility=${arg.facility}&device=${arg.device}`)
    .then((response) => response.data as typeof sample);
};
