import { withRunnerLayout } from "@/layouts/RunnerLayout/RunnerLayout";
import { RunnerPageComponent } from "@/page-components/RunnerPageComponent/RunnerPageComponent";
import { JSX } from "react";
//import {GetServerSideProps} from "next";
//import {EmployeesStore} from "@/stores/EmployeesStore";

function RunnerPage({ }): JSX.Element {
  
	return (
		<>
			<RunnerPageComponent />
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
  
export default withRunnerLayout(RunnerPage);