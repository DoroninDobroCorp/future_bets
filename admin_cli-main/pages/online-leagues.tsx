import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { OnlineLeaguesPageComponent } from "@/page-components/OnlineLeaguesPageComponent/OnlineLeaguesPageComponent";

function OnlineLeaguesPage({ }): JSX.Element {
    return (
        <>
            <OnlineLeaguesPageComponent />
        </>
    )
}

export default withLayout(OnlineLeaguesPage);