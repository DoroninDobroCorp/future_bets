import {withEmployeesLayout} from "@/layouts/EmployeesLayout/EmployeesLayout";
import { EmployeesPageComponent } from "@/page-components/EmployeesPageComponent/EmployeesPageComponent";
import { JSX } from "react";
//import {GetServerSideProps} from "next";
//import {EmployeesStore} from "@/stores/EmployeesStore";

function EmployeesPage({ }): JSX.Element {

    return (
        <>
            <EmployeesPageComponent />
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

export default withEmployeesLayout(EmployeesPage);