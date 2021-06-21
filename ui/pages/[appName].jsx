import { useRouter } from "next/router";
import { useQuery } from "react-query";
import { useMemo } from "react";
import CssBaseline from "@material-ui/core/CssBaseline";
import MaUTable from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import { useTable, useFlexLayout } from "react-table";
import TableContainer from "@material-ui/core/TableContainer";
import Paper from "@material-ui/core/Paper";
import { css } from "@emotion/react";
import Link from "next/link";

async function fetchStuff(url) {
  try {
    const response = await fetch(url);
    return response.json();
  } catch (err) {
    throw new Error("Network response was not ok");
  }
}

const App = () => {
  const router = useRouter();
  const { appName } = router.query;

  const hierarchy = [
    "global",
    "app-global",
    "environment-class",
    "app-environment-class",
    "environment",
    "app",
  ];

  const { status, data, error, isFetching } = useQuery(["stuff", appName], () =>
    fetchStuff(`http://localhost:3000/test-data/${appName}.json`)
  );

  const {
    status: appFetchStatus,
    data: apps,
    error: appFetchError,
    isFetching: isFetchingApps,
  } = useQuery(["apps"], () =>
    fetchStuff(`http://localhost:3000/test-data/applications.json`)
  );

  if (
    status === "loading" ||
    appFetchStatus === "loading" ||
    isFetching ||
    isFetchingApps
  ) {
    return <span>Loading...</span>;
  }

  if (status === "error" || appFetchStatus === "error") {
    if (status === "error") {
      return <span>Error: {error.message}</span>;
    }
    if (appFetchStatus === "error") {
      return <span>Error: {appFetchError.message}</span>;
    }
    return (
      <>
        <span>Error: {error.message}</span>
        <span>Error: {appFetchError.message}</span>
      </>
    );
  }

  return (
    <div id="top">
      <ul
        css={css`
          list-style: none;
          display: flex;
          flex-wrap: wrap;
        `}
      >
        {apps.map((app) => (
          <Link href={`/${app}`} key={app}>
            <li
              css={css`
                border: 1px solid black;
                margin: 1rem;
                padding: 1rem;
              `}
            >
              <a href={`/${app}`}>{app}</a>
            </li>
          </Link>
        ))}
      </ul>
      ;
      <Table
        processed={getProcessed(data, hierarchy)}
        hierarchy={hierarchy}
        appName={appName}
      />
    </div>
  );
};

const Table = (props) => {
  const { processed, hierarchy, appName } = props;

  const columns = useMemo(
    () => [
      {
        Header: "Property",
        accessor: "property",
        minWidth: 150,
        maxWidth: 250,
        width: 250,
      },
      ...makeColumns(processed),
    ],
    [processed]
  );

  const env = useMemo(
    () => makeEnv(processed, hierarchy),
    [processed, hierarchy]
  );

  return (
    <>
      <CssBaseline />
      {hierarchy.map((hier) => (
        <>
          <SubTable
            columns={columns}
            key={`${appName}-${hier}`}
            data={env[hier]}
            hier={hier}
          />
        </>
      ))}
    </>
  );
};

function SubTable({ columns, data, hier }) {
  // Use the state and functions returned from useTable to build your UI
  const { getTableProps, headerGroups, rows, prepareRow } = useTable(
    {
      columns,
      data,
    },
    useFlexLayout
  );

  // Render the UI for your table
  return (
    <>
      <h1>{hier}</h1>
      <TableContainer
        component={Paper}
        className="scroll-table"
        key={`${hier}-scroll-table`}
      >
        <MaUTable {...getTableProps()} aria-label="sticky table" size="small">
          <TableHead>
            {headerGroups.map((headerGroup) => (
              <>
                <TableRow
                  {...headerGroup.getHeaderGroupProps([
                    {
                      style: {
                        fontWeight: "bold",
                        left: 0,
                        top: 0,
                        position: "sticky",
                        zIndex: 2,
                        backgroundColor: "lightgray",
                      },
                    },
                  ])}
                >
                  {headerGroup.headers.map((column) => (
                    <TableCell
                      {...column.getHeaderProps([
                        {
                          style: {
                            fontWeight: "bold",
                          },
                        },
                      ])}
                    >
                      {column.render("Header")}
                    </TableCell>
                  ))}
                </TableRow>
              </>
            ))}
          </TableHead>
          <TableBody>
            {rows.map((row) => {
              prepareRow(row);
              return (
                <>
                  <TableRow {...row.getRowProps()}>
                    {row.cells.map((cell, i) => (
                      <>
                        {i === 0 ? (
                          <TableCell
                            component="th"
                            scope="row"
                            {...cell.getCellProps([
                              {
                                style: {
                                  fontWeight: "bold",
                                  wordBreak: "break-all",
                                },
                              },
                            ])}
                          >
                            {cell.render("Cell")}
                          </TableCell>
                        ) : (
                          <TableCell
                            {...cell.getCellProps([
                              {
                                style: {
                                  wordBreak: "break-all",
                                },
                              },
                            ])}
                          >
                            {cell.render("Cell")}
                          </TableCell>
                        )}
                      </>
                    ))}
                  </TableRow>
                </>
              );
            })}
          </TableBody>
        </MaUTable>
      </TableContainer>
    </>
  );
}

function makeColumns(processed) {
  return Object.keys(processed.environments).map((env) => ({
    Header: env,
    accessor: env,
    minWidth: 150,
    maxWidth: 250,
    width: 250,
  }));
}

export default App;

function getProcessed(data, hierarchy) {
  return Object.entries(data[0]).reduce(
    (acc, [, value]) => {
      const { Global, SubResource, Environment, Objects, Properties } = value;

      Object.keys(Properties).forEach((propKey) => {
        acc.keys[propKey] = null;
      });

      if (Global === false && SubResource === false) {
        acc.environments[Environment] = {};
      }

      hierarchy.forEach((val) => {
        if (acc.environments[Environment] !== undefined) {
          if (Objects != null) {
            const hObjects = Objects.filter(
              (item) => item.Hierarchy.hierarchyName === val
            );
            if (hObjects != null && hObjects.length > 0) {
              acc.environments[Environment][val] = hObjects[0].Item.data;
            }
          }
        }
      });

      if (Global === true) {
        acc.globals[Environment] = {};
      }

      hierarchy.forEach((val) => {
        if (acc.globals[Environment] !== undefined) {
          if (Objects != null) {
            const hObjects = Objects.filter(
              (item) => item.Hierarchy.hierarchyName === val
            );
            if (hObjects != null && hObjects.length > 0) {
              acc.globals[Environment][val] = hObjects[0].Item.data;
            }
          }
        }
      });

      return acc;
    },
    { environments: {}, globals: {}, keys: {} }
  );
}

function makeEnv(processed, hierarchy) {
  return hierarchy.reduce((acc, hier) => {
    acc[hier] = [];
    Object.keys(processed.keys).forEach((key) => {
      if (hier === "global" || hier === "app-global") {
        const keyOrNull =
          processed.globals.all != null &&
          processed.globals.all[hier] != null &&
          processed.globals.all[hier][key] != null
            ? processed.globals.all[hier][key]
            : null;

        const row = Object.keys(processed.environments).reduce((a, env) => {
          a[env] = keyOrNull;
          return a;
        }, {});

        row.property = key;
        acc[hier][acc[hier].length] = row;
      }

      if (hier === "app" || hier === "environment") {
        const row = Object.keys(processed.environments).reduce((a, env) => {
          const keyOrNull =
            processed.environments != null &&
            processed.environments[env] != null &&
            processed.environments[env][hier] != null &&
            processed.environments[env][hier][key] != null
              ? processed.environments[env][hier][key]
              : null;
          a[env] = keyOrNull;
          return a;
        }, {});
        row.property = key;
        acc[hier][acc[hier].length] = row;
      }
    });
    return acc;
  }, {});
}
