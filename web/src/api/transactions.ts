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

const sampleOATransactionResponse = {
  data: [
    {
      id: "b2e3dc44-c0b5-4c9d-a8d1-2fb3ace40fff",
      businessTransactionId: "2880ba20-1fe2-4bce-bc5d-a642e3687f5a",
      lpn: "UVW2345",
      customerid: "70948456",
      jobid: "jid",
      facility: "1230",
      device: "101",
      extra: {
        error: "TEST",
        leaveAt: "2023-11-27T10:21:17.844838Z",
        steps: "payment_exit_error",
      },
      entryLane: "101",
      exitLane: "201",
      createdAt: "2023-11-27T10:18:11.805379Z",
      updatedAt: "2023-11-27T10:18:11.805379Z",
    },
  ],
  metadata: {
    page: 0,
    perPage: 10,
    totalData: 1,
    totalPage: 1,
  },
};

export interface OATransactionsQuery {
  startAt?: Date;
  endAt?: Date;
  exitLane?: string;
  entryLane?: string;
  lpn?: string;
  facility?: string;
  jobId?: string;
  page?: number;
  perPage?: number;
}

export const getOATransactions = async (query: OATransactionsQuery) => {
  return axios
    .get(`/transactions/oa`, { params: { ...query } })
    .then((response) => response.data as typeof sampleOATransactionResponse);
};

const sampleIntegratorTransactionResponse = {
  data: [
    {
      businessTransactionId: "45525bff-59d1-4777-a47b-c158f09fc52b",
      lpn: "UVW2345",
      integratorId: "",
      status: "success",
      amount: "10",
      error: "",
      integratorName: "TNG",
      extra: {},
      taxData: {},
      createdAt: "2023-11-30T16:24:39.483222Z",
      updatedAt: "2023-11-30T16:24:39.483222Z",
    },
  ],
  metadata: {
    page: 0,
    perPage: 50,
    totalData: 15,
    totalPage: 1,
  },
};

export interface IntegratorTransactionsQuery {
  startAt?: Date;
  endAt?: Date;
  lpn?: string;
  status?: string;
  integratorName?: string;
  page?: number;
  perPage?: number;
}

export const getIntegratorTransactions = async (
  query: IntegratorTransactionsQuery,
) => {
  return axios
    .get(`/transactions/integrator`, { params: { ...query } })
    .then(
      (response) => response.data as typeof sampleIntegratorTransactionResponse,
    );
};
