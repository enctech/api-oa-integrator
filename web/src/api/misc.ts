import axios from "./axios";
import moment from "moment/moment";

interface MiscResponse {
  totalIn: number;
  totalOut: number;
  totalPayment: number;
}

interface IntegratorStatusResponse {
  integrators: { integrator: string; status: string }[];
  snb: { facility: string; status: string }[];
}

export const misc = async () => {
  return axios
    .get("/misc/", {
      params: {
        startAt: moment().startOf("day").utc().toDate(),
        endAt: moment().endOf("day").utc().toDate(),
      },
    })
    .then((res) => res.data as MiscResponse);
};

export const integratorStatus = async () => {
  return axios
    .get("/misc/integrator")
    .then((res) => res.data as IntegratorStatusResponse);
};
