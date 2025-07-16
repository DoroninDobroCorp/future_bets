import Document, { Html, Head, Main, NextScript, DocumentContext, DocumentInitialProps } from "next/document";
import { JSX } from "react";

class MyDocument extends Document {

  static async getInitialProps(ctx: DocumentContext): Promise<DocumentInitialProps> {
      const initialProps = await Document.getInitialProps(ctx);
      return { ...initialProps };
  }

  render(): JSX.Element {
      return (
          <Html lang="ru">
              <Head>
                  <link rel="preconnect" href="https://fonts.gstatic.com" />
                  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Noto+Sans+KR:wght@300;400;500;700&display=swap" />
              </Head>
              <body>
                  <Main />
                  <NextScript />
              </body>
          </Html>
      );
  }
}

export default MyDocument;
