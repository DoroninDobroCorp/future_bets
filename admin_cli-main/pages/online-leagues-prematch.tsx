import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { OnlineLeaguesPrematchPageComponent } from "@/page-components/OnlineLeaguesPrematchPageComponent/OnlineLeaguesPrematchPageComponent";

function OnlineLeaguesPrematchPage({ }): JSX.Element {
    return (
        <>
            <OnlineLeaguesPrematchPageComponent />
        </>
    )
}

export default withLayout(OnlineLeaguesPrematchPage);