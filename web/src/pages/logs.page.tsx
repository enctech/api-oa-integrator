import React, { useEffect, useRef } from "react";
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TablePagination,
  TableRow,
  TextField,
} from "@mui/material";
import { useQuery } from "react-query";
import { getOALogs } from "../api/transactions";
import { JsonViewer } from "@textea/json-viewer";
import { ClearIcon, DateTimePicker } from "@mui/x-date-pickers";
import dayjs, { Dayjs } from "dayjs";
import { useSearchParams } from "react-router-dom";
import { Controller, SubmitHandler, useForm } from "react-hook-form";
import moment from "moment/moment";
import IconButton from "@mui/material/IconButton";
import { useDebounce } from "@uidotdev/usehooks";

interface FormData {
  startAt?: Dayjs | null;
  endAt?: Dayjs | null;
  message: string;
  field: string;
}

const LogsPage = () => {
  const perPagesDefault = useRef([100, 500, 1000]);
  const [currentQueryParameters, setSearchParams] = useSearchParams();
  const newParams = useRef(new URLSearchParams());

  const debouncedSearchTerm = useDebounce(currentQueryParameters, 300);

  const { control, register, handleSubmit, watch, setValue } =
    useForm<FormData>({
      defaultValues: {
        message: currentQueryParameters.get("message") || "",
        field: currentQueryParameters.get("field") || "",
        startAt: dayjs(
          moment(currentQueryParameters.get("startAt")).local().toDate(),
        ),
        endAt: dayjs(
          moment(currentQueryParameters.get("endAt")).local().toDate(),
        ),
      },
    });

  watch(["startAt", "endAt", "message", "field"]);

  const handleFieldChange = () => {
    handleSubmit(onSubmit)();
  };

  const onSubmit: SubmitHandler<FormData> = (formData) => {
    if (formData.startAt && formData.startAt.isValid()) {
      newParams.current.set("startAt", formData.startAt.toISOString());
    }
    if (formData.endAt && formData.endAt.isValid()) {
      newParams.current.set("endAt", formData.endAt.toISOString());
    }

    newParams.current.set("message", formData.message || "");
    newParams.current.set("field", formData.field || "");
    newParams.current.set("page", "0");

    setSearchParams(newParams.current);
  };

  const { data, refetch } = useQuery(
    "getOALogs",
    () =>
      getOALogs({
        startAt: currentQueryParameters.get("startAt")
          ? moment(currentQueryParameters.get("startAt")).utc().toDate()
          : undefined,
        endAt: currentQueryParameters.get("endAt")
          ? moment(currentQueryParameters.get("endAt")).utc().toDate()
          : undefined,
        message: currentQueryParameters.get("message") || undefined,
        field: currentQueryParameters.get("field") || undefined,
        page: +(currentQueryParameters.get("page") || "0"),
        perPage: +(
          currentQueryParameters.get("perPage") ||
          `${perPagesDefault.current[0]}`
        ),
      }),
    {
      enabled: false, // Temporarily disabled - query is too slow, see PR #2
      refetchInterval: 5000,
    },
  );

  useEffect(() => {
    if (debouncedSearchTerm) {
      refetch();
    }
  }, [debouncedSearchTerm]);

  useEffect(() => {
    newParams.current.set("page", "0");
    newParams.current.set("perPage", "100");
    if (currentQueryParameters.get("startAt")) {
      newParams.current.set(
        "startAt",
        moment(currentQueryParameters.get("startAt")).utc().toISOString(),
      );
    }
    if (currentQueryParameters.get("endAt")) {
      newParams.current.set(
        "endAt",
        moment(currentQueryParameters.get("endAt")).utc().toISOString(),
      );
    }
    newParams.current.set(
      "message",
      currentQueryParameters.get("message") || "",
    );
    newParams.current.set("field", currentQueryParameters.get("field") || "");

    setSearchParams(newParams.current);
  }, []);

  return (
    <div className="p-4 px-8">
      <div className={"flex"}>
        <div className={"py-4"}>
          <Controller
            name="startAt"
            control={control}
            render={({ field }) => (
              <div className="flex">
                <DateTimePicker
                  format={"DD/MM/YYYY hh:mm"}
                  className="flex-1 w-full"
                  {...field}
                  onChange={(_) => {}}
                  label="Start Date"
                  onAccept={(value) => {
                    if (!value) return;
                    field.onChange(value);
                    handleFieldChange();
                  }}
                />
                <IconButton
                  onClick={() => {
                    setValue("startAt", null);
                    newParams.current.delete("startAt");
                    setSearchParams(newParams.current);
                  }}
                >
                  <ClearIcon />
                </IconButton>
              </div>
            )}
          />
        </div>
        <div className={"p-4"}>
          <Controller
            name="endAt"
            control={control}
            render={({ field }) => (
              <div className="flex">
                <DateTimePicker
                  slotProps={{
                    toolbar: {
                      toolbarFormat: "YYYY",
                      toolbarPlaceholder: "??",
                    },
                    actionBar: {
                      actions: ["clear", "accept"],
                    },
                  }}
                  className="flex-1 w-full"
                  {...field}
                  format={"DD/MM/YYYY hh:mm"}
                  label="End Date"
                  onAccept={(value) => {
                    if (!value) return;
                    field.onChange(value);
                    handleFieldChange();
                  }}
                />
                <IconButton
                  onClick={() => {
                    setValue("endAt", null);
                    newParams.current.delete("endAt");
                    setSearchParams(newParams.current);
                  }}
                >
                  <ClearIcon />
                </IconButton>
              </div>
            )}
          />
        </div>
        <div className={"p-4"}>
          <TextField
            label="Message"
            {...register("message", {
              onChange: (_) => handleFieldChange(),
            })}
          />
        </div>
        <div className={"p-4"}>
          <TextField
            label="Field"
            {...register("field", {
              onChange: (_) => handleFieldChange(),
            })}
          />
        </div>
      </div>
      <TableContainer component={Paper} className="mt-4">
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Created At</TableCell>
              <TableCell>Level</TableCell>
              <TableCell>Message</TableCell>
              <TableCell>Fields</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data?.data?.map((row) => (
              <TableRow key={row.id}>
                <TableCell>
                  <div className="w-[5rem]">
                    {moment(row.createdAt)
                      .local()
                      .format("DD/MM/yyyy hh:mm:ss.SSSSSS A")}
                  </div>
                </TableCell>
                <TableCell>
                  <div className="w-[5rem]">{row.level}</div>
                </TableCell>
                <TableCell>
                  <div className="w-[20rem]">{row.message}</div>
                </TableCell>
                <TableCell>
                  <JsonViewer
                    className="w-[20rem]"
                    value={row.fields}
                    defaultInspectControl={(_, __) => false}
                    collapseStringsAfterLength={false}
                  />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={perPagesDefault.current}
          component="div"
          count={data?.metadata.totalData || 0}
          rowsPerPage={parseInt(
            currentQueryParameters.get("perPage") ||
              `${perPagesDefault.current[0]}`,
            10,
          )}
          page={parseInt(currentQueryParameters.get("page") || "0", 10)}
          onPageChange={(_: any, newPage: number) => {
            newParams.current.set("page", newPage.toString());
            setSearchParams(newParams.current);
          }}
          onRowsPerPageChange={(event: { target: { value: string } }) => {
            newParams.current.set("perPage", event.target.value);
            setSearchParams(newParams.current);
          }}
        />
      </TableContainer>
    </div>
  );
};

export default LogsPage;
