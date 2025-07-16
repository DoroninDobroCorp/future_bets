import { withClientLayout } from "@/layouts/ClientLayout/ClientLayout";
import { ClientPageComponent } from "@/page-components/ClientPageComponent/ClientPageComponent";
import { JSX } from "react";
//import {GetServerSideProps} from "next";
//import {EmployeesStore} from "@/stores/EmployeesStore";

function ClientPage({ }): JSX.Element {
	return (
		<>
			<ClientPageComponent userId={1} />
		</>
	)
  }

/*export const getServerSideProps: GetServerSideProps = async (context) => {
    const { token } = context.query;

    if (!token) {
        return {
            notFound: true,
        };
    }

    const userId = await EmployeesStore.sendToken(token as string);

    if (userId === 401) {
        return {
            notFound: true
        };
    }

    return {
        props: {
            userId: userId
        },
    };
};*/


export default withClientLayout(ClientPage);