import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { AllLeaguesPageComponent } from "@/page-components/AllLeaguesPageComponent/AllLeaguesPageComponent";

function AllLeaguesPage({ }): JSX.Element {
    return (
        <>
            <AllLeaguesPageComponent />
        </>
    )
}

export default withLayout(AllLeaguesPage);