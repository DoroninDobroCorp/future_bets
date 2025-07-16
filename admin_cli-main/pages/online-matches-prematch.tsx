import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { OnlineMatchesPrematchPageComponent } from "@/page-components/OnlineMatchesPrematchPageComponent/OnlineMatchesPrematchPageComponent";

function OnlineMatchesPrematchPage({ }): JSX.Element {

    return (
        <>
            <OnlineMatchesPrematchPageComponent />
        </>
    )
}

export default withLayout(OnlineMatchesPrematchPage);