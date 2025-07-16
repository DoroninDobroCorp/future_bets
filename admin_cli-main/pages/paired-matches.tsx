import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { PairedMatchesPageComponent } from "@/page-components/PairedMatchesPageComponent/PairedMatchesPageComponent";

function PairedMatchesPage({ }): JSX.Element {
  
	return (
		<>
			<PairedMatchesPageComponent />
		</>
	)
  }
  
  export default withLayout(PairedMatchesPage);