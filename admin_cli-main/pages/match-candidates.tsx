import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { MatchCandidatesPageComponent } from "@/page-components/MatchCandidatesPageComponent/MatchCandidatesPageComponent";

function MatchCandidatesPage({ }): JSX.Element {
  
	return (
		<>
			<MatchCandidatesPageComponent />
		</>
	)
  }
  
  export default withLayout(MatchCandidatesPage);