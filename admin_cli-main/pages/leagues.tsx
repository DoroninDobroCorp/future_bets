import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { LeaguesPageComponent } from "@/page-components/LeaguesPageComponent/LeaguesPageComponent";
//import {GetServerSideProps} from "next";
//import {EmployeesStore} from "@/stores/EmployeesStore";

function LeaguesPage({ }): JSX.Element {
	return (
		<>
			<LeaguesPageComponent />
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
  
export default withLayout(LeaguesPage);