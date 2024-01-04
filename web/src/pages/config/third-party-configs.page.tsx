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
import { getIntegratorConfigs, IntegratorConfigs } from "../../api/config";
import { useNavigate } from "react-router-dom";
import { AdminOnly } from "../../components/auth-guard";

const ThirdPartyConfigsPage = () => {
  const navigate = useNavigate();

  const perPagesDefault = useRef([100, 500, 1000]);

  const [rowsPerPage, setRowsPerPage] = useState(perPagesDefault.current[0]);
  const [page, setPage] = useState(0);

  const { data, refetch } = useQuery("getIntegratorConfigs", () =>
    getIntegratorConfigs(),
  );

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
    navigate(`/3rd-party-configs/${id}`);
  };

  const createNewConfig = () => {
    navigate(`/3rd-party-configs/new`);
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
              <TableCell>Provider ID</TableCell>
              <TableCell>Client ID</TableCell>
              <TableCell>Service Provider ID</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data?.map((row) => (
              <IntegratorConfig
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

const IntegratorConfig = ({
  row,
  handleRowClick,
}: {
  row: IntegratorConfigs;
  handleRowClick: (id: string) => void;
}) => {
  return (
    <TableRow
      className={"cursor-pointer"}
      key={row.id}
      onClick={() => handleRowClick(row.id!)}
    >
      <TableCell>{row.name}</TableCell>
      <TableCell>{row.providerId}</TableCell>
      <TableCell>{row.clientId}</TableCell>
      <TableCell>{row.serviceProviderId}</TableCell>
    </TableRow>
  );
};

export default ThirdPartyConfigsPage;
