import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { UnpairedMatchesPageComponent } from "@/page-components/UnpairedMatchesPageComponent/UnpairedMatchesPageComponent";

function UnpairedMatchesPage({ }): JSX.Element {
  
	return (
		<>
			<UnpairedMatchesPageComponent />
		</>
	)
  }
  
  export default withLayout(UnpairedMatchesPage);