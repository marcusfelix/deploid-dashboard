interface Task {
    title: string;
    body: string;
    comments: number;
    issue: number;
    labels: string[];
    state: "open" | "closed";
    user: {
        name: string;
        login: string;
        photo: string;
    };
    created: string;
}

interface Comment {
    body: string;
    user: {
        name: string;
        login: string;
        photo: string;
    };
    created: string;
}

interface Commit {
    sha: string;
    message: string;
    author: {
      name: string;
      avatar: string;
    };
    updated: Date;
  }
  
  interface Branch {
    name: string;
    sha: string;
  }

  interface Workflow {
    branch: string;
    name: string;
    status: string;
    created: Date;
  }