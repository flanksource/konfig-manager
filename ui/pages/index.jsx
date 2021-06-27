import Head from "next/head";
import Link from "next/link";
import { css } from "@emotion/react";
import { useQuery } from "react-query";
import styles from "../styles/Home.module.css";

async function fetchStuff(url) {
  try {
    const response = await fetch(url);
    return response.json();
  } catch (err) {
    throw new Error("Network response was not ok");
  }
}

export default function Home() {
  const {
    status,
    data: apps,
    error,
    isFetching,
  } = useQuery(
    ["apps"],
    () => {
      if (window.DEMO_MODE) {
        return fetchStuff(
          `${process.env.serverPrefix}/test-data/applications.json`
        );
      }
      return fetchStuff(`${process.env.serverPrefix}/api/applications`);
    },
    { staleTime: 1000, refetchOnWindowFocus: false }
  );

  if (status === "loading" || isFetching) {
    return <span>Loading...</span>;
  }

  if (status === "error") {
    return <span>Error: {error.message}</span>;
  }

  return (
    <div className={styles.container}>
      <Head>
        <title>Konfig Manager</title>
        <meta name="description" content="KonfigManager" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className={styles.main}>
        {Array.isArray(apps) ? (
          <ul
            css={css`
              list-style: none;
              display: flex;
              flex-wrap: wrap;
            `}
          >
            {apps.map((app) => (
              <li
                key={app}
                css={css`
                  border: 1px solid black;
                  margin: 1rem;
                  padding: 1rem;
                `}
              >
                <Link href={`/${app}`}>
                  <a href={`/${app}`}>{app}</a>
                </Link>
              </li>
            ))}
          </ul>
        ) : (
          <div>Application list not found</div>
        )}
      </main>
    </div>
  );
}
