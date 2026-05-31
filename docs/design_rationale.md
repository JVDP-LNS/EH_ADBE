# Design Rationale
I have doing a lot of agentic automation work in my internship so I thought maybe I could apply it to EH Modding. 



## Basic requirements 
Sorted roughly in order of importance to me:
- Fully Free (Both development and usage)
- Agentic mod dev features
- Manual mod dev features
- Cross Platform
- Lightweight
- Efficient



## Basic components
- A frontend or UI
- A backend for handling logic and orchestration
- An AI model



## High level Design
Frontend and Backend design is simple enough and implementation is easy using coding assistants. The main architectural hurdle is getting a working AI LLM for every user. The problem is that running AI models requires lots of compute (GPU) which most EH players probably dont have much of (Assuming most are young mobile players). 

The solution is something I randomly came across on google. Kaggle (or Google Colab) provide a online environment with compute included (with certain limits) for free. The intended use case is to aid in Data Science and AI research. For the purposes of ADBE, the use case is technically "Generation of synthetic data", which can be classified under Data Science. 

### Kaggle limits
- 30 hrs GPU compute per week
- At most 12 hr per notebook session
- either P100 or 2x T4 GPUs 
    - T4 GPUs give total 32 GB VRAM, with inference speed about same as RTX 3050

The normal way to use kaggle is to use their website and interactively run the code. However, they also provide APIs (and keys) where we can push code to be run as a notebook. There are still some limitations to this approach (needs investigation).

### Connecting Local and Kaggle
Tunneling is an option but Kaggle will probably not approve of it. Instead, we will setup an outbound websocket connection from kaggle notebook to our server. We can use localhost.run (SSH based) to get a public URL for our server and inject that into the kaggle notebook.
