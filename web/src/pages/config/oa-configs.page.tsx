import React, { useEffect, useRef, useState } from "react";
import {
  Button,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TablePagination,
  TableRow,
} from "@mui/material";
import { useQuery } from "react-query";
import { getOAConfigs, OAConfigResponse } from "../../api/config";
import { useNavigate } from "react-router-dom";
import { getOAHealth } from "../../api/health";
import { AdminOnly } from "../../components/auth-guard";

const OAConfigsPage = () => {
  const navigate = useNavigate();

  const perPagesDefault = useRef([100, 500, 1000]);

  const [rowsPerPage, setRowsPerPage] = useState(perPagesDefault.current[0]);
  const [page, setPage] = useState(0);

  const { data, refetch } = useQuery("getOAConfigs", () => getOAConfigs());

  useEffect(() => {
    refetch().then();
  }, []);

  const handleChangePage = (_: any, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: { target: { value: string } }) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const handleRowClick = (id: string) => {
    navigate(`/oa-configs/${id}`);
  };

  const createNewConfig = () => {
    navigate(`/oa-configs/new`);
  };

  return (
    <div>
      <AdminOnly>
        <Button
          variant="contained"
          sx={{
            backgroundColor: "#3f3100",
            color: "#fff0bf",
            "&:hover": {
              backgroundColor: "#fff0bf",
              color: "#3f3100",
            },
          }}
          onClick={createNewConfig}
        >
          New
        </Button>
      </AdminOnly>
      <TableContainer component={Paper} className="mt-4">
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Facilities</TableCell>
              <TableCell>Devices</TableCell>
              <TableCell>Status</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data?.map((row) => (
              <OAConfigRow
                key={row.id}
                row={row}
                handleRowClick={handleRowClick}
              />
            ))}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={perPagesDefault.current}
          component="div"
          count={data?.length || 0}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={handleChangePage}
          onRowsPerPageChange={handleChangeRowsPerPage}
        />
      </TableContainer>
    </div>
  );
};

const OAConfigRow = ({
  row,
  handleRowClick,
}: {
  row: OAConfigResponse;
  handleRowClick: (id: string) => void;
}) => {
  const { data: oaHealth } = useQuery(
    "getOAHealth",
    () =>
      row.devices.length > 0 && row.facilities.length > 0
        ? getOAHealth({
            device: row.devices[0],
            facility: row.facilities[0],
          })
        : Promise.resolve(null),
    {},
  );

  return (
    <TableRow
      className={"cursor-pointer"}
      key={row.id}
      onClick={() => handleRowClick(row.id)}
    >
      <TableCell>{row.name}</TableCell>
      <TableCell>{row.facilities?.join(", ")}</TableCell>
      <TableCell>{row.devices?.join(", ")}</TableCell>
      <TableCell>
        {oaHealth?.oa === "up" ? (
          <div
            className="w-5 h-5 rounded-full
                inline-flex items-center justify-center
                bg-green-500"
          ></div>
        ) : (
          <div
            className="w-5 h-5 rounded-full
                inline-flex items-center justify-center
                bg-red-500"
          ></div>
        )}
      </TableCell>
    </TableRow>
  );
};

export default OAConfigsPage;
