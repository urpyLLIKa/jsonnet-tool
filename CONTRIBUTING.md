## Developer Certificate of Origin + License

By contributing to GitLab B.V., You accept and agree to the following terms and
conditions for Your present and future Contributions submitted to GitLab B.V.
Except for the license granted herein to GitLab B.V. and recipients of software
distributed by GitLab B.V., You reserve all right, title, and interest in and to
Your Contributions. All Contributions are subject to the following DCO + License
terms.

[DCO + License](https://gitlab.com/gitlab-org/dco/blob/master/README.md)

All Documentation content that resides under the [doc/ directory](/doc) of this
repository is licensed under Creative Commons:
[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).

_This notice should stay as the first item in the CONTRIBUTING.md file._

### Merge requests contributors workflow

The instructions below assume that you have access to the repository at its canonical location on GitLab.com. If you are onboarded as GitLab Inc. team-member, the access is granted during onboarding. The canonical repository does not run production workloads, as the changes are mirrored to a separate instance that connects to production environments.

The workflow to make a merge request is as follows:

1. Clone the project and create a feature branch.
1. Make the necessary changes
1. Commit and submit the MR with the changes.
1. Ensure that MR description contains links to relevant resources, **and** explain why the specific change is being made.
1. Apply the `Contribution` label, as well as any other applicable labels (Service labels and similar).
1. If the request is a part of corrective action for an active incident, assign the MR to `SRE on-call`. Current on-call can be found in the [production channel](https://gitlab.slack.com/archives/C101F3796), in the `sre-oncall` user group.

### Merge request reviewers workflow

As a reviewer, you need to ensure a certain level of quality for the MR you're assigned to review.
Keep the following flow in mind:

1. Ensure that you are familiar with the [contributing guidelines](/CONTRIBUTING.md#contributing-guidelines).
1. In order to make matters simpler, assume that the contributor has a limited perspective into how services run, so double check the intention of the MR.
1. Ensure that the MR description has context on why the change is made, and links to the applicable resources such as issues, epics and other related MRs. Valid descriptions allow others to understand context even when they are not participating in the work. It also increases long-term readability in case the MR needs to be referenced in future.
1. Ensure that the MR has labels. Labels make it simpler to track multiple changes over time.
1. Review the CI pipeline comments from `ops-gitlab-net` user, as they contain the full pipeline run from the operational instance. The pipelines in the MR widget only run syntax checks.
1. Before merging the MR, ensure that the `Reviewer Check-list` section in the MR description is addressed.
1. Once the changes are reviewed, merge the MR and ensure that the changes are successfully applied to all applicable environments. Not doing so carries a risk of a failed rollout which does block regular operations such as deployments.

## Project workflow

The workflow below is generally used by the members of the Infrastructure department for changes more advanced than a configuration value change.

### Project mirroring

This project has two locations:

1. Canonical project on GitLab.com: The project is public, accessible to everyone in [the gitlab-com namespace](https://gitlab.com/gitlab-com).
1. Mirror project on [ops.gitlab.net](https://ops.gitlab.net): The project is on the Ops instance, accessible only to the Infrastructure department team members.

Detailed explanation can be found on [the Infrastructure project classification handbook page](https://about.gitlab.com/handbook/engineering/infrastructure/projects/), but in short:

1. Canonical project enables everyone to contribute, without gaining direct access to production environments. Not everyone has production access, but everyone should be able to propose at minimum simple configuration changes. The canonical project has no production access.
1. Mirror project is connected to production environments and contains sensitive information. The project is also used in cases when GitLab.com is not fully operational, and a change needs to be authored and applied.
1. Mirroring is one-way, from Canonical to Mirror project, unless there is an emergency. See [emergency workflow](#workflow-during-emergencies)

Do note a couple of additional items:

1. This workflow might not be as convenient as granting everyone access to the complete project.
1. Mirroring between the Canonical and Mirror project can take more than a few seconds, which puts an extra load on the reviewer to ensure the validity of the change.

However, this workflow enables us to work within the company values. We work as transparently and collaboratively as possible, while ensuring we are aligned with the compliance requirements.

To further support the two repository approach, we utilize the following configuration:

1. Canonical project uses push mirroring.
1. Mirror project uses the following supporting CI variables
    * `MIRROR` - CI configuration used to mark the project as a mirror.
    * `GITLAB_API_TOKEN` - Read access token by `ops-gitlab-net` user on GitLab.com, to send CI pipeline statuses the Canonical project.

### Workflow during emergencies

In emergencies such as GitLab.com incident, or a bug that prevents regular workflow, MRs on the mirror project on ops.gitlab.net can be used.

The MR's on the Mirror project have 3 approvals required to merge. Merging in the Mirror project would break the one-way mirroring, so as a safe guard we want and additional review. This is ok, during such emergencies.

Once the emergency is resolved, one of the 3 approvers need to manually sync the changes between the Mirror and the Cannonical project.
This can be done in the local copy:

1. Pull the changes from the Canonical project.
1. Pull the changes from the Mirror project.
1. Resolve any possible conflicts.
1. Push the changes to Canonical project.

After the last step, mirroring should return to normal. Confirm this by pushing a feature branch on the Canonical project, and verify that the branch has been mirrored on the Mirror project.
