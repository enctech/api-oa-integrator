import axios from "./axios";
import moment from "moment/moment";

interface MiscResponse {
  integrators: { integrator: string; status: string }[];
  snb: { facility: string; status: string }[];
  totalIn: number;
  totalOut: number;
  totalPayment: Map<string, string>;
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
