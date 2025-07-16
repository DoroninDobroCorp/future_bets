import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { OnlineMatchesPageComponent } from "@/page-components/OnlineMatchesPageComponent/OnlineMatchesPageComponent";

function OnlineMatchesPage({ }): JSX.Element {

    return (
        <>
            <OnlineMatchesPageComponent />
        </>
    )
}

export default withLayout(OnlineMatchesPage);