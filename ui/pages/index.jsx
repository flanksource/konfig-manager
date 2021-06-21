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
  } = useQuery(["apps"], () =>
    fetchStuff(`http://localhost:3000/test-data/applications.json`)
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
        <ul
          css={css`
            list-style: none;
            display: flex;
            flex-wrap: wrap;
          `}
        >
          {apps.map((app) => (
            <Link href={`/${app}`}>
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
      </main>
    </div>
  );
}
