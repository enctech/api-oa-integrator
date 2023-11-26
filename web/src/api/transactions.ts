import axios from "./axios";
import moment from "moment";

const sample = {
  data: [
    {
      id: "ab752b66-ef35-4264-85aa-dc21b7f07c09",
      level: "error",
      message: "failed migrate up no change",
      fields: {
        timestamp: "2023-11-20T22:27:13.240+0800",
      },
      createdAt: "2023-11-20T14:27:13.24Z",
    },
  ],
  metadata: {
    page: 1,
    perPage: 10,
    totalData: 118,
    totalPage: 12,
  },
};

export interface OALogsQuery {
  before: Date;
  after: Date;
  message: string | null;
  field: string | null;
  page: number;
  perPage: number;
}

export const getOALogs = async (query: OALogsQuery) => {
  return axios
    .get(
      `/transactions/logs?before=${moment
        .utc(query.before)
        .toISOString()}&after=${moment
        .utc(query.after)
        .toISOString()}&message=${query.message}&fields=${query.field}&page=${
        query.page
      }&perPage=${query.perPage}`,
    )
    .then((response) => response.data as typeof sample);
};
