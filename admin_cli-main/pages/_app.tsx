import "@/styles/globals.css";
import type { AppProps } from "next/app";
import Head from "next/head";
import { JSX } from "react";
import { ToastContainer } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';

export default function App({ Component, pageProps }: AppProps): JSX.Element {
  return <>
    <Head>
      <title>Livebets AdminCli</title>
      <meta name="description" content="livebets" />
      <meta name="og:title" content="Livebets Admin Client" />
      <meta name="og:description" content="livebets" />
      <meta name="og:locale" content="ru_RU" />
      <meta name="og:type" content="website" />
      <link rel="icon" href="/favicon.ico" />
    </Head>
    <Component {...pageProps} />
    <ToastContainer  position='bottom-right' autoClose={400}/>
  </>;
}
