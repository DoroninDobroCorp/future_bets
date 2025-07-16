import { withLayout } from "@/layouts/MainLayout/MainLayout";
import { JSX } from "react";
import { LeagueCandidatesPageComponent } from "@/page-components/LeagueCandidatesPageComponent/LeagueCandidatesPageComponent";

function LeagueCandidatesPage({ }): JSX.Element {
  
	return (
		<>
			<LeagueCandidatesPageComponent />
		</>
	)
  }
  
  export default withLayout(LeagueCandidatesPage);