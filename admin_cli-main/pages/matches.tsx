import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { MatchesPageComponent } from "@/page-components/MatchesPageComponent/MatchesPageComponent";

function MatchesPage({ }): JSX.Element {
  
	return (
		<>
			<MatchesPageComponent />
		</>
	)
  }
  
  export default withLayout(MatchesPage);