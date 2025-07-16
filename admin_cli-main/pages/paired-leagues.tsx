import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { PairedLeaguesPageComponent } from "@/page-components/PairedLeaguesPageComponent/PairedLeaguesPageComponent";

function PairedLeaguesPage({ }): JSX.Element {
  
	return (
		<>
			<PairedLeaguesPageComponent />
		</>
	)
  }
  
  export default withLayout(PairedLeaguesPage);